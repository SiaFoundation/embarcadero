export type SwapSatusLocal = 'creatingANewSwap' | 'loadingAnExistingSwap'

export type SwapStatusRemoteInt = 0 | 1 | 2 | 3 | 4

type SwapStatusRemote =
  | 'waitingForYouToAccept'
  | 'waitingForCounterpartyToAccept'
  | 'waitingForCounterpartyToFinish'
  | 'waitingForYouToFinish'
  | 'transactionComplete'

export type SwapStatus = SwapSatusLocal | SwapStatusRemote

const swapStatusRemoteMapping = {
  0: 'waitingForYouToAccept',
  1: 'waitingForCounterpartyToAccept',
  2: 'waitingForCounterpartyToFinish',
  3: 'waitingForYouToFinish',
  4: 'transactionComplete',
} as Record<SwapStatusRemoteInt, SwapStatusRemote>

export function getSwapStatusRemote(
  status?: SwapStatusRemoteInt
): SwapStatusRemote | undefined {
  if (status === undefined) {
    return undefined
  }
  return swapStatusRemoteMapping[status]
}
