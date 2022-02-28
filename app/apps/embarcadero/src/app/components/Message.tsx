import { ExclamationTriangleIcon, InfoCircledIcon } from '@radix-ui/react-icons'
import { Alert, Box, Flex, Text } from '@siafoundation/design-system'

type Props = {
  variant?: React.ComponentProps<typeof Alert>['variant']
  message: string
}

export function Message({ variant, message }: Props) {
  return (
    <Alert variant={variant} css={{ width: '100%' }}>
      <Flex gap="2">
        <Box css={{ color: '$hiContrast', position: 'relative', top: '4px' }}>
          {variant === 'red' ? (
            <ExclamationTriangleIcon />
          ) : (
            <InfoCircledIcon />
          )}
        </Box>
        <Text css={{ lineHeight: '20px' }}>{message}</Text>
      </Flex>
    </Alert>
  )
}
