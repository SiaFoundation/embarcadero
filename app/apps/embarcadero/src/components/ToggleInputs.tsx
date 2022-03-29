import { ArrowDown16, Box, Flex } from '@siafoundation/design-system'

type Props = {
  onToggle?: () => void
  disabled?: boolean
}

export function ToggleInputs({ onToggle, disabled }: Props) {
  return (
    <Box css={{ height: '$2', zIndex: 1, order: 2 }}>
      <Box
        tabIndex={2}
        onClick={onToggle}
        css={{
          position: 'relative',
          top: '-15px',
          height: '40px',
          width: '40px',
          backgroundColor: '$loContrast',
          borderRadius: '15px',
        }}
      >
        <Flex
          align="center"
          justify="center"
          css={{
            backgroundColor: '$brandGray5',
            borderRadius: '10px',
            position: 'absolute',
            transform: 'translate(-50%, -50%)',
            left: '50%',
            color: '$hiContrast',
            top: '50%',
            height: '30px',
            width: '30px',
            transition: 'background 0.1s linear',
            '&:hover': !disabled && {
              cursor: 'pointer',
              backgroundColor: '$brandGray7',
            },
          }}
        >
          <ArrowDown16 />
        </Flex>
      </Box>
    </Box>
  )
}
