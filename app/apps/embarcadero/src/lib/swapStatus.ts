export type SwapStatusLocal = 'creatingANewSwap' | 'loadingAnExistingSwap'

export type SwapStageRemoteInt = 0 | 1 | 2 | 3 | 4

export type SwapStatusRemote =
  | 'waitingForYouToAccept'
  | 'waitingForCounterpartyToAccept'
  | 'waitingForCounterpartyToFinish'
  | 'waitingForYouToFinish'
  | 'transactionComplete'

export type SwapStatus = SwapStatusLocal | SwapStatusRemote

const swapStatusRemoteMapping = {
  0: 'waitingForYouToAccept',
  1: 'waitingForCounterpartyToAccept',
  2: 'waitingForCounterpartyToFinish',
  3: 'waitingForYouToFinish',
  4: 'transactionComplete',
} as Record<SwapStageRemoteInt, SwapStatusRemote>

export function getSwapStatusRemote(
  status?: SwapStageRemoteInt
): SwapStatusRemote | undefined {
  if (status === undefined) {
    return undefined
  }
  return swapStatusRemoteMapping[status]
}

export const localSwapStatuses = ['creatingANewSwap', 'loadingAnExistingSwap']
