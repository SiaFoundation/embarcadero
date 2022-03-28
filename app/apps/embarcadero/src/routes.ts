import { SwapStatus } from './lib/swapStatus'

export const routes = {
  home: '/',
  create: '/create',
  input: '/input',
  swap: '/swap',
}

export const swapStatusToRoute: Record<SwapStatus, keyof typeof routes> = {
  createANewSwap: 'create',
  openASwap: 'input',
  waitingForYouToAccept: 'swap',
  waitingForCounterpartyToAccept: 'swap',
  waitingForYouToFinish: 'swap',
  waitingForCounterpartyToFinish: 'swap',
  swapTransactionPending: 'swap',
  swapTransactionConfirmed: 'swap',
}
