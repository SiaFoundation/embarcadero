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
import { useMemo } from 'react'
import { useHistory } from 'react-router-dom'
import useSWR from 'swr'
import { routes } from '../routes'
import { usePathParams } from './useHashParam'

export type SwapTransaction = {
  siacoinInputs: SiacoinInput[]
  siafundInputs: SiafundInput[]
  siacoinOutputs: SiacoinOutput[]
  siafundOutputs: SiafundOutput[]
  signatures: TransactionSignature[]
}

export type SwapSatusLocal = 'creatingANewSwap' | 'loadingAnExistingSwap'

export type SwapStatusRemoteInt = 0 | 1 | 2 | 3

type SwapStatusRemote =
  | 'waitingForYouToAccept'
  | 'waitingForCounterpartyToAccept'
  | 'waitingForCounterpartyToFinish'
  | 'waitingForYouToFinish'

const swapStatusRemoteMapping = {
  0: 'waitingForYouToAccept',
  1: 'waitingForCounterpartyToAccept',
  2: 'waitingForCounterpartyToFinish',
  3: 'waitingForYouToFinish',
} as Record<SwapStatusRemoteInt, SwapStatusRemote>

function getSwapStatusRemote(
  status?: SwapStatusRemoteInt
): SwapStatusRemote | undefined {
  if (status === undefined) {
    return undefined
  }
  return swapStatusRemoteMapping[status]
}

export type SwapStatus = SwapSatusLocal | SwapStatusRemote

type SwapSummary = {
  receiveSF: boolean
  receiveSC: boolean
  payFee: boolean
  amountSC: string
  amountSF: string
  amountFee: string
  stage: SwapStatusRemoteInt
}

type Response = {
  summary: SwapSummary
  swap: SwapTransaction
}

export function useSwap(hash?: string) {
  const { route } = usePathParams()
  const response = useSWR<Response>(hash, async () => {
    const res = await axios({
      method: 'post',
      url: 'http://localhost:8080/api/summarize',
      headers: {
        'Content-Type': 'application/json',
      },
      data: {
        raw: hash,
      },
    })
    return res.data
  })

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

  let localStatus: SwapSatusLocal | undefined = undefined
  if (route === routes.create.slice(1)) {
    localStatus = 'creatingANewSwap'
  } else if (route === routes.input.slice(1)) {
    localStatus = 'loadingAnExistingSwap'
  }

  return {
    isValidating: response.isValidating,
    summary: response.data?.summary,
    transaction: response.data?.swap,
    status: getSwapStatusRemote(response.data?.summary.stage) || localStatus,
    offerSc,
    sf,
    sc,
  }
}
