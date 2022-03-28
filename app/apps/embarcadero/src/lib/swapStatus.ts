export type SwapStatusLocal = 'createANewSwap' | 'openASwap'

export type SwapStageRemoteInt = 1 | 2 | 3 | 4 | 5 | 6

export type SwapStatusRemote =
  | 'waitingForYouToAccept'
  | 'waitingForCounterpartyToAccept'
  | 'waitingForYouToFinish'
  | 'waitingForCounterpartyToFinish'
  | 'swapTransactionPending'
  | 'swapTransactionConfirmed'

export type SwapStatus = SwapStatusLocal | SwapStatusRemote

const swapStatusRemoteMapping = {
  1: 'waitingForYouToAccept',
  2: 'waitingForCounterpartyToAccept',
  3: 'waitingForYouToFinish',
  4: 'waitingForCounterpartyToFinish',
  5: 'swapTransactionPending',
  6: 'swapTransactionConfirmed',
} as Record<SwapStageRemoteInt, SwapStatusRemote>

export function getSwapStatusRemote(
  status?: SwapStageRemoteInt
): SwapStatusRemote | undefined {
  if (status === undefined) {
    return undefined
  }
  return swapStatusRemoteMapping[status]
}

export const localSwapStatuses = ['createANewSwap', 'openASwap']
