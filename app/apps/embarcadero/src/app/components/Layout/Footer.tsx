import {
  Flex,
  Panel,
  Separator,
  Text,
  Tooltip,
} from '@siafoundation/design-system'
import { Fragment } from 'react'
import { useConnectivity } from '../../hooks/useConnectivity'
import { useConsensus } from '../../hooks/useConsensus'
import { useSiaStats } from '../../hooks/useSiaStats'
import { useWallet } from '../../hooks/useWallet'
import { NetworkStatus } from '../NetworkStatus'

export function Footer() {
  const { siad } = useConnectivity()
  const { data: siaStats } = useSiaStats()
  const { data: consensus, error: errorC } = useConsensus()
  const { data: wallet } = useWallet()

  const haveData = consensus && siaStats && wallet

  const isSynced = haveData && consensus.height >= siaStats.block_height

  const color = errorC ? 'red' : isSynced ? 'green' : 'yellow'

  return (
    <Panel
      css={{
        position: 'fixed',
        bottom: '$3',
        right: '$5',
        padding: '$2 $3',
      }}
    >
      <Flex gap="2" align="center">
        {siad && (
          <Fragment>
            <Tooltip content="Current transaction fee">
              <Text size="1">
                {((Number(wallet?.dustthreshold) / Math.pow(10, 24)) * 1024) /
                  0.001}{' '}
                mS / KB
              </Text>
            </Tooltip>
            <Separator orientation="vertical" />
            <Tooltip
              content={
                isSynced
                  ? 'Block height'
                  : `Block height: ${consensus?.height} / ${siaStats?.block_height}`
              }
            >
              <Text size="1">{consensus?.height}</Text>
            </Tooltip>
            <Separator orientation="vertical" />
          </Fragment>
        )}
        <NetworkStatus
          variant={color}
          content={!siad ? 'Disconnected' : isSynced ? 'Synced' : 'Syncing'}
        />
      </Flex>
    </Panel>
  )
}