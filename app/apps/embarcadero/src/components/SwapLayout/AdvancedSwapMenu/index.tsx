import {
  Flex,
  Heading,
  IconButton,
  Popover,
  PopoverContent,
  PopoverTrigger,
  Settings24,
} from '@siafoundation/design-system'
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
          <Heading>Advanced</Heading>
          <Details />
        </Flex>
      </PopoverContent>
    </Popover>
  )
}
