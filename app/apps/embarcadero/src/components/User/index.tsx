import { Flex } from '@siafoundation/design-system'
import { UserMenu } from './UserMenu'
import { Wallet } from './Wallet'

export function User() {
  return (
    <Flex gap="2" align="center">
      <Wallet />
      <UserMenu size="2" />
    </Flex>
  )
}
