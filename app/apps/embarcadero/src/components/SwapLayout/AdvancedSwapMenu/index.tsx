import {
  Flex,
  Heading,
  IconButton,
  Popover,
  PopoverContent,
  PopoverTrigger,
  Settings16,
} from '@siafoundation/design-system'
import { Details } from './Details'

export function AdvancedSwapMenu() {
  return (
    <Popover>
      <PopoverTrigger asChild>
        <IconButton css={{ transform: 'scale(1.5)' }}>
          <Settings16 />
        </IconButton>
      </PopoverTrigger>
      <PopoverContent align="end" css={{ padding: '$3', minWidth: '400px' }}>
        <Flex direction="column" gap="3">
          <Heading>Advanced</Heading>
          <Details />
        </Flex>
      </PopoverContent>
    </Popover>
  )
}
