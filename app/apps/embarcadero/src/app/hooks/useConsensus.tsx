import useSWR from 'swr'
import { getApi } from '../../config'
import { ConsensusGET } from '@siafoundation/sia-js'
import { handleResponse } from '../lib/handleResponse'

export function useConsensus() {
  return useSWR<ConsensusGET>(
    'consensus',
    async () => {
      const response = await fetch(getApi('/api/consensus'))
      return handleResponse(response)
    },
    {
      refreshInterval: 10_000,
      errorRetryInterval: 10_000,
    }
  )
}
