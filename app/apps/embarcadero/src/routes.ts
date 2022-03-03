import { SwapStatus } from './lib/swapStatus'

export const routes = {
  home: '/',
  create: '/create',
  input: '/input',
  swap: '/swap',
}

export const swapStatusToRoute: Record<SwapStatus, keyof typeof routes> = {
  creatingANewSwap: 'create',
  loadingAnExistingSwap: 'input',
  waitingForYouToAccept: 'swap',
  waitingForCounterpartyToAccept: 'swap',
  waitingForCounterpartyToFinish: 'swap',
  waitingForYouToFinish: 'swap',
  transactionComplete: 'swap',
}
