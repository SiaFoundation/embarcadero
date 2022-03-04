import { toSiacoins } from '@siafoundation/sia-js'
import { useMemo } from 'react'
import { useWallet } from '@siafoundation/sia-react'
import { api } from '../config'
import BigNumber from 'bignumber.js'

type Props = {
  isOffer?: boolean
  currency: 'SF' | 'SC'
  value?: BigNumber
}

export function useHasBalance({ currency, isOffer, value }: Props) {
  const { data: wallet } = useWallet({
    api,
  })

  return useMemo(() => {
    if (!isOffer || !value) {
      return true
    }
    if (currency === 'SC') {
      return toSiacoins(wallet?.confirmedsiacoinbalance || 0).gte(value)
    }
    const sfBalance = new BigNumber(wallet?.siafundbalance || 0)

    return sfBalance >= value
  }, [isOffer, currency, value, wallet])
}
