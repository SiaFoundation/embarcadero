import {
  Flex,
  Heading,
  IconButton,
  Popover,
  PopoverContent,
  PopoverTrigger,
  Settings24,
} from '@siafoundation/design-system'
import { CopyTxnId } from '../../CopyTxnId'
import { Details } from './Details'

export function AdvancedSwapMenu() {
  return (
    <Popover>
      <PopoverTrigger asChild>
        <IconButton size="2">
          <Settings24 />
        </IconButton>
      </PopoverTrigger>
      <PopoverContent align="end" css={{ padding: '$2', minWidth: '400px' }}>
        <Flex direction="column" gap="2">
          <Flex align="center" justify="between">
            <Heading>Advanced</Heading>
            <CopyTxnId />
          </Flex>
          <Details />
        </Flex>
      </PopoverContent>
    </Popover>
  )
}
