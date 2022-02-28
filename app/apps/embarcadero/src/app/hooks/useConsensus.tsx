import useSWR from 'swr'
import { getApi } from '../../config'
import { ConsensusGET } from '@siafoundation/sia-js'

interface SWRError extends Error {
  info?: string
  status?: number
}

export function useConsensus() {
  return useSWR<ConsensusGET>(
    'consensus',
    async () => {
      let res: Response | undefined = undefined
      // try {
      res = await fetch(getApi('/api/consensus'))
      // } catch (e) {
      //   if (e instanceof Error) {
      //     throw e
      //   }
      // }

      // If the status code is not in the range 200-299,
      // we still try to parse and throw it.
      if (!res || !res.ok) {
        const message = await res?.text()
        const error: SWRError = new Error(
          message || 'An error occurred while fetching the data.'
        )
        // Attach extra info to the error object.
        error.status = res?.status || 500
        throw error
      } else {
        return res.json()
      }
    },
    {
      refreshInterval: 10_000,
      errorRetryInterval: 10_000,
    }
  )
}
