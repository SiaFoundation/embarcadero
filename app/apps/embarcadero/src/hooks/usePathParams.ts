import { useMemo } from 'react'
import { useLocation } from 'react-router-dom'

export function usePathParams() {
  const location = useLocation()
  const [_, route, encodedPath] = location.pathname.split('/')
  const path = useMemo(
    () => (encodedPath ? decodeURIComponent(encodedPath) : undefined),
    [encodedPath]
  )

  return {
    route,
    path,
  }
}
