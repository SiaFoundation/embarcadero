import {
  Alert,
  Box,
  Flex,
  Information16,
  Text,
  Warning16,
} from '@siafoundation/design-system'

type Props = {
  icon?: React.ReactNode
  variant?: React.ComponentProps<typeof Alert>['variant']
  message: string
}

export function Message({ variant, icon, message }: Props) {
  return (
    <Alert variant={variant} css={{ width: '100%' }}>
      <Flex gap="2">
        <Box css={{ color: '$hiContrast', position: 'relative', top: '4px' }}>
          {icon || (variant === 'red' ? <Warning16 /> : <Information16 />)}
        </Box>
        <Text css={{ lineHeight: '20px' }}>{message}</Text>
      </Flex>
    </Alert>
  )
}
