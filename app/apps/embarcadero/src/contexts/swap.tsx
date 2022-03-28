import React, {
  createContext,
  useMemo,
  useContext,
  useCallback,
  useState,
  useEffect,
  useRef,
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
import { routes } from '../routes'
import { usePathParams } from '../hooks/usePathParams'
import { downloadJsonFile } from '../lib/download'
import { useHistory } from 'react-router-dom'
import {
  getSwapStatusRemote,
  SwapStatusLocal,
  SwapStatus,
  SwapStageRemoteInt,
} from '../lib/swapStatus'
import { swapTxnSchema } from '../lib/validate'
import { api } from '../config'

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
}

type State = {
  id?: string
  txn?: SwapTransaction
  summary?: SwapSummary
  status?: SwapStatus
  offerSc: boolean
  sf?: BigNumber
  sc?: BigNumber
  isValidating: boolean
  loadTxnFromFile: (file: File) => void
  resetTxn: () => void
  downloadTxn: () => void
  loadTxn: (txn: SwapTransaction) => void
  signTxn: (step: 'accept' | 'finish') => void
  fileReadError?: string
  txnError?: string
}

const SwapContext = createContext({} as State)
export const useSwap = () => useContext(SwapContext)

type Props = {
  children: React.ReactNode
}

export function SwapProvider({ children }: Props) {
  const { route: currentRoute } = usePathParams()
  const history = useHistory()

  const [txn, setTxn] = useState<SwapTransaction>()
  const [id, setId] = useState<string>()
  const [summary, setSummary] = useState<SwapSummary>()
  const [isValidating, setIsValidating] = useState<boolean>(false)

  const [fileReadError, setFileReadError] = useState<string>()
  const [txnError, setTxnError] = useState<string>()

  const resetTxn = useCallback(() => {
    setId(undefined)
    setSummary(undefined)
    setTxn(undefined)
    setTxnError(undefined)
  }, [setId, setSummary, setTxn, setTxnError])

  const loadTxn = useCallback(
    (txn: SwapTransaction) => {
      const func = async () => {
        try {
          setIsValidating(true)
          const { data, error } = await fetchSummary(txn)

          if (data) {
            const { id, summary } = data
            setFileReadError(undefined)
            setId(id)
            setSummary(summary)
            setTxn(txn)
            console.log(id, summary, txn)
            const nextRoute: keyof typeof routes = 'swap'
            if (currentRoute !== nextRoute) {
              history.push(routes.swap)
            }
          } else {
            setTxnError('Error fetching transaction summary.')
          }
        } catch (e) {
          setFileReadError('Invalid transaction file.')
        } finally {
          setIsValidating(false)
        }
      }
      func()
    },
    [setTxn, history, setFileReadError, currentRoute]
  )

  const validateAndLoadTxnFile = useCallback(
    (fileData: string) => {
      const func = async () => {
        try {
          setIsValidating(true)

          const txn = JSON.parse(fileData) as SwapTransaction

          const validated = swapTxnSchema.validate(txn)
          if (validated.error) {
            console.log(validated.error)
            setFileReadError('Invalid transaction file.')
            return
          }

          loadTxn(txn)
        } catch (e) {
          console.log(e)
          setFileReadError('Invalid transaction file.')
        } finally {
          setIsValidating(false)
        }
      }
      func()
    },
    [loadTxn, setFileReadError]
  )

  const loadTxnFromFile = useCallback(
    (file: File) => {
      const reader = new FileReader()
      const decoder = new TextDecoder()

      reader.onabort = () => {
        console.log('file reading was aborted')
      }
      reader.onerror = () => {
        setFileReadError('Failed to read transaction file.')
      }
      reader.onload = () => {
        const bin = reader.result
        const str = decoder.decode(bin as ArrayBuffer)

        if (str) {
          validateAndLoadTxnFile(str)
        } else {
          setFileReadError('Empty transaction file.')
        }
      }
      reader.readAsArrayBuffer(file)
    },
    [setFileReadError, validateAndLoadTxnFile]
  )

  const signTxn = useCallback(
    (step: 'accept' | 'finish') => {
      const func = async () => {
        try {
          const response = await axios({
            method: 'post',
            url: `${api}/api/${step}`,
            headers: {
              'Content-Type': 'application/json',
            },
            data: {
              swap: txn,
            },
          })

          loadTxn(response.data.swap)
        } catch (e) {
          if (e instanceof Error) {
            setTxnError('Error signing transaction.')
          }
        }
      }
      func()
    },
    [txn, setTxnError, loadTxn]
  )

  const downloadTxn = useCallback(() => {
    if (id && txn) {
      downloadJsonFile(`embc_txn_${id.slice(0, 6)}`, txn)
    }
  }, [id, txn])

  const { sc, sf, offerSc } = useMemo(() => {
    if (!summary) {
      return {
        sc: undefined,
        sf: undefined,
        offerSc: false,
      }
    }

    const { amountSC, amountSF, receiveSF } = summary

    const sc = toSiacoins(new BigNumber(amountSC))
    const sf = new BigNumber(amountSF)
    const offerSc = !!receiveSF

    return {
      sc,
      sf,
      offerSc,
    }
  }, [summary])

  let localStatus: SwapStatusLocal | undefined = undefined
  if (currentRoute === routes.create.slice(1)) {
    localStatus = 'createANewSwap'
  } else if (currentRoute === routes.input.slice(1)) {
    localStatus = 'openASwap'
  }

  const status = getSwapStatusRemote(summary?.stage) || localStatus

  const ref = useRef({
    txn,
  })

  useEffect(() => {
    ref.current.txn = txn
  }, [txn])

  useEffect(() => {
    if (
      txn &&
      status &&
      (
        [
          'waitingForCounterpartyToFinish',
          'swapTransactionPending',
        ] as SwapStatus[]
      ).includes(status)
    ) {
      startPolling(() => {
        if (ref.current.txn) {
          loadTxn(ref.current?.txn)
        }
      })
    } else {
      stopPolling()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status])

  const value: State = {
    id,
    isValidating,
    summary,
    txn,
    status,
    offerSc,
    sf,
    sc,
    loadTxnFromFile,
    resetTxn,
    downloadTxn,
    loadTxn,
    signTxn,
    fileReadError,
    txnError,
  }

  return <SwapContext.Provider value={value}>{children}</SwapContext.Provider>
}

async function fetchSummary(txn: SwapTransaction) {
  try {
    const res = await axios({
      method: 'post',
      url: `${api}/api/summarize`,
      headers: {
        'Content-Type': 'application/json',
      },
      data: {
        swap: txn,
      },
    })
    return {
      data: res.data as SummarizeResponse,
      error: undefined,
    }
  } catch (e) {
    return {
      error: e,
      data: undefined,
    }
  }
}

let interval: NodeJS.Timer | null = null

function startPolling(func: () => void): void {
  if (!interval) {
    interval = setInterval(func, 5_000)
  }
}

function stopPolling(): void {
  if (interval) {
    clearInterval(interval)
    interval = null
  }
}
