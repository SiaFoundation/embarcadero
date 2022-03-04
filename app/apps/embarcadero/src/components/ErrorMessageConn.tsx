import { Message } from './Message'
import { Connectivity, useConnectivity } from '../hooks/useConnectivity'

export function ErrorMessageConn() {
  const connectivity = useConnectivity()

  const connError = getConnError(connectivity)

  if (!connError) {
    return null
  }

  return <Message variant="red" message={connError} />
}

function getConnError(conn: Connectivity) {
  if (!conn.embd) {
    return 'Connect to embd to continue'
  }
  if (!conn.siad) {
    return 'Connect to siad to continue'
  }
  if (!conn.wallet) {
    return 'Unlock wallet to continue'
  }
  return ''
}
