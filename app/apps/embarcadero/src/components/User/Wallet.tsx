import { Flex, Panel, Separator, Text } from '@siafoundation/design-system'
import { useWallet } from '@siafoundation/sia-react'
import { api } from '../../config'

export function Wallet() {
  const { data: wallet } = useWallet({
    api,
  })

  if (!wallet?.unlocked) {
    return null
  }

  return (
    <Panel>
      <Flex gap="2" align="center" css={{ height: '$6', padding: '0 $2' }}>
        <Text css={{ fontWeight: '600' }}>
          {(
            Number(wallet?.confirmedsiacoinbalance || 0) / Math.pow(10, 24)
          ).toLocaleString()}{' '}
          SC
        </Text>
        <Separator orientation="vertical" />
        <Text css={{ fontWeight: '600' }}>
          {Number(wallet?.siafundbalance || 0).toLocaleString()} SF
        </Text>
      </Flex>
    </Panel>
  )
}
