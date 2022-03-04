import { useSwap } from '../contexts/swap'
import { useHasBalance } from './useHasBalance'

export function useTxnHasBalance() {
  const { offerSc, sc, sf } = useSwap()

  return useHasBalance({
    value: offerSc ? sc : sf,
    isOffer: true,
    currency: offerSc ? 'SC' : 'SF',
  })
}
