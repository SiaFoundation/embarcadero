import {
  Flex,
  Panel,
  Separator,
  Text,
  Tooltip,
} from '@siafoundation/design-system'
import { Fragment } from 'react'
import { useConnectivity } from '../../hooks/useConnectivity'
import {
  useWallet,
  useConsensus,
  useSiaStatsNetworkStatus,
} from '@siafoundation/sia-react'
import { NetworkStatus } from '../NetworkStatus'
import { api } from '../../config'
import { useSettings } from '../../hooks/useSettings'

export function Footer() {
  const { siad } = useConnectivity()
  const { settings } = useSettings()
  const { data: siaStats } = useSiaStatsNetworkStatus({
    disabled: !settings.siaStats,
  })
  const { data: consensus, error: errorC } = useConsensus({
    api,
  })
  const { data: wallet } = useWallet({
    api,
  })

  const isSynced = consensus?.synced

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
              <Text size="1" css={{ fontFamily: '$mono', lineHeight: '1' }}>
                {((Number(wallet?.dustthreshold) / Math.pow(10, 24)) * 1024) /
                  0.001}{' '}
                mS / KB
              </Text>
            </Tooltip>
            <Separator orientation="vertical" />
            <Tooltip
              content={
                settings.siaStats
                  ? `Block height: ${consensus?.height} / ${siaStats?.block_height}`
                  : 'Block height'
              }
            >
              <Text size="1" css={{ fontFamily: '$mono' }}>
                {consensus?.height}
              </Text>
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
