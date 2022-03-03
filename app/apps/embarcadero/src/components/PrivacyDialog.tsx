import {
  Alert,
  Box,
  DialogTitle,
  Information16,
  Text,
  Flex,
  Heading,
  Switch,
  Paragraph,
} from '@siafoundation/design-system'
import { useSettings } from '../hooks/useSettings'

export function PrivacyDialog() {
  const { settings, setSettings } = useSettings()

  return (
    <Box>
      <Flex direction="column" gap="4">
        <DialogTitle>Privacy</DialogTitle>
        <Paragraph size="1">
          This app uses the following third-party APIs, all external APIs are
          not required and can be toggled on or off.
        </Paragraph>
        <Alert>
          <Flex direction="column" gap="2">
            <Flex gap="1" align="center">
              <Box>
                <Information16 />
              </Box>
              <Heading size="1" css={{ flex: 1 }}>
                SiaStats
              </Heading>
              <Switch
                checked={settings.siaStats}
                onCheckedChange={(val) =>
                  setSettings((settings) => ({ ...settings, siaStats: val }))
                }
              />
            </Flex>
            <Paragraph size="1">
              The app fetches price information and block height from the
              SiaStats API.
            </Paragraph>
          </Flex>
        </Alert>
      </Flex>
    </Box>
  )
}
