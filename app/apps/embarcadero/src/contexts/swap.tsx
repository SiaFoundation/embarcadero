import React, {
  createContext,
  useMemo,
  useContext,
  useCallback,
  useState,
} from 'react'
import {
  SiacoinInput,
  SiacoinOutput,
  SiafundInput,
  SiafundOutput,
  toSiacoins,
  TransactionSignature,
} from '@siafoundation/sia-js'
import axios from 'axios'
import BigNumber from 'bignumber.js'
import useSWR from 'swr'
import { routes } from '../routes'
import { usePathParams } from '../hooks/usePathParams'
import { downloadFile } from '../lib/download'
import { useHistory } from 'react-router-dom'
import {
  getSwapStatusRemote,
  SwapStatusLocal,
  SwapStatus,
  SwapStageRemoteInt,
} from '../lib/swapStatus'

export type SwapTransaction = {
  siacoinInputs: SiacoinInput[]
  siafundInputs: SiafundInput[]
  siacoinOutputs: SiacoinOutput[]
  siafundOutputs: SiafundOutput[]
  signatures: TransactionSignature[]
}

type SwapSummary = {
  receiveSF: boolean
  receiveSC: boolean
  payFee: boolean
  amountSC: string
  amountSF: string
  amountFee: string
  stage: SwapStageRemoteInt
}

type SummarizeResponse = {
  id: string
  summary: SwapSummary
  swap: SwapTransaction
}

type State = {
  id?: string
  isValidating: boolean
  summary?: SwapSummary
  transaction?: SwapTransaction
  status?: SwapStatus
  offerSc: boolean
  sf?: BigNumber
  sc?: BigNumber
  raw?: string
  loadTransactionFromFile: (file: File) => void
  clearTransaction: () => void
  downloadTransaction: () => void
  setTransaction: (raw: string) => void
  signTransaction: (step: 'accept' | 'finish') => void
  fileReadError?: string
  transactionError?: string
}

const SwapContext = createContext({} as State)
export const useSwap = () => useContext(SwapContext)

type Props = {
  children: React.ReactNode
}

export function SwapProvider({ children }: Props) {
  const [raw, setRaw] = useState<string>()
  const [fileReadError, setFileReadError] = useState<string>()
  const [transactionError, setTransactionError] = useState<string>()
  const { route: currentRoute } = usePathParams()
  const history = useHistory()

  const validateTransactionFile = useCallback(
    (data: string) => {
      // TODO: add validation
      if (data) {
        setFileReadError(undefined)
        setRaw(data)
        const nextRoute: keyof typeof routes = 'swap'
        if (currentRoute !== nextRoute) {
          history.push(routes.swap)
        }
      } else {
        setFileReadError('Invalid transaction file')
      }
    },
    [setRaw, history, setFileReadError, currentRoute]
  )

  const loadTransactionFromFile = useCallback(
    (file: File) => {
      const reader = new FileReader()
      const decoder = new TextDecoder()

      reader.onabort = () => {
        console.log('file reading was aborted')
      }
      reader.onerror = () => {
        setFileReadError('Failed to read transaction file')
      }
      reader.onload = () => {
        const bin = reader.result
        const str = decoder.decode(bin as ArrayBuffer)
        validateTransactionFile(str)
      }
      reader.readAsArrayBuffer(file)
    },
    [setFileReadError, validateTransactionFile]
  )

  const setTransaction = useCallback(
    (raw: string) => {
      setRaw(raw)
    },
    [setRaw]
  )

  const clearTransaction = useCallback(() => {
    setRaw(undefined)
  }, [setRaw])

  const signTransaction = useCallback(
    (step: 'accept' | 'finish') => {
      const func = async () => {
        try {
          const response = await axios({
            method: 'post',
            url: `http://localhost:8080/api/${step}`,
            headers: {
              'Content-Type': 'application/json',
            },
            data: {
              raw,
            },
          })

          setTransaction(response.data.raw)
        } catch (e) {
          if (e instanceof Error) {
            setTransactionError(e.message)
          }
        }
      }
      func()
    },
    [raw, setTransactionError, setTransaction]
  )

  const { route } = usePathParams()
  const response = useSWR<SummarizeResponse>(raw, async () => {
    const res = await axios({
      method: 'post',
      url: 'http://localhost:8080/api/summarize',
      headers: {
        'Content-Type': 'application/json',
      },
      data: {
        raw,
      },
    })
    return res.data
  })

  const id = response.data?.id

  const downloadTransaction = useCallback(() => {
    if (id && raw) {
      downloadFile(`transaction_${id.slice(0, 5)}`, raw)
    }
  }, [id, raw])

  const offerSc = useMemo(() => {
    return !!response.data?.summary.receiveSF
  }, [response])

  const sc = useMemo(() => {
    if (!response.data?.summary) {
      return undefined
    }

    const { amountSC } = response.data.summary

    return toSiacoins(new BigNumber(amountSC))
  }, [response])

  const sf = useMemo(() => {
    if (!response.data?.summary) {
      return undefined
    }

    const { amountSF } = response.data.summary

    return new BigNumber(amountSF)
  }, [response])

  let localStatus: SwapStatusLocal | undefined = undefined
  if (route === routes.create.slice(1)) {
    localStatus = 'creatingANewSwap'
  } else if (route === routes.input.slice(1)) {
    localStatus = 'loadingAnExistingSwap'
  }

  const value: State = {
    id,
    isValidating: response.isValidating,
    summary: response.data?.summary,
    transaction: response.data?.swap,
    status: getSwapStatusRemote(response.data?.summary.stage) || localStatus,
    offerSc,
    sf,
    sc,
    raw,
    loadTransactionFromFile,
    clearTransaction,
    downloadTransaction,
    setTransaction,
    signTransaction,
    fileReadError,
    transactionError,
  }
  console.log(id, !!raw, response.isValidating)

  return <SwapContext.Provider value={value}>{children}</SwapContext.Provider>
}
