import useLocalStorageState from 'use-local-storage-state'

const defaultConfig = {
  siaStats: true,
}

export function useSettings() {
  const [settings, setSettings] = useLocalStorageState('v0/settings', {
    ssr: false,
    defaultValue: defaultConfig,
  })

  return { settings, setSettings }
}
