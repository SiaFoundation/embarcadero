export type SwapStatusLocal = 'createANewSwap' | 'openASwap'

export type SwapStatusRemote =
  | 'waitingForYouToAccept'
  | 'waitingForCounterpartyToAccept'
  | 'waitingForYouToFinish'
  | 'waitingForCounterpartyToFinish'
  | 'swapTransactionPending'
  | 'swapTransactionConfirmed'

export type SwapStatus = SwapStatusLocal | SwapStatusRemote

export const localSwapStatuses = ['createANewSwap', 'openASwap']
