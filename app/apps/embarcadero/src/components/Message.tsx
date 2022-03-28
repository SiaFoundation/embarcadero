import {
  Alert,
  Box,
  Flex,
  Information16,
  NextOutline16,
  CheckmarkOutline16,
  Text,
  Warning16,
} from '@siafoundation/design-system'

type Props = {
  variant?: 'step' | 'info' | 'error' | 'success'
  message: string
}

export function Message({ variant = 'step', message }: Props) {
  return (
    <Alert
      variant={variant === 'error' ? 'red' : 'gray'}
      css={{ width: '100%' }}
    >
      <Flex gap="1">
        <Box css={{ color: '$hiContrast', position: 'relative', top: '1px' }}>
          {variant === 'step' && <NextOutline16 />}
          {variant === 'info' && <Information16 />}
          {variant === 'error' && <Warning16 />}
          {variant === 'success' && <CheckmarkOutline16 />}
        </Box>
        <Text css={{ lineHeight: '20px' }}>{message}</Text>
      </Flex>
    </Alert>
  )
}
