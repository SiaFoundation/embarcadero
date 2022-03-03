import { useWallet } from '@siafoundation/sia-react'
import { api } from '../config'

export function useConnectivity() {
  const w = useWallet({
    api,
  })

  // TODO: remove any after updating package
  const connError = w.error

  // Any error fetching wallet data means siad is not connected
  const siad = !connError

  // 500 to wallet is a siad issue, any other error means embd is not connected
  const embd = !(connError && connError.status !== 500)

  const wallet = !!w.data?.unlocked

  return {
    all: siad && embd && wallet,
    connections: siad && embd,
    siad,
    embd,
    wallet,
  }
}

export type Connectivity = ReturnType<typeof useConnectivity>
