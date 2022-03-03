import { ArrowDownIcon } from '@radix-ui/react-icons'
import {
  Box,
  Button,
  Flex,
  triggerErrorToast,
} from '@siafoundation/design-system'
import axios from 'axios'
import { useCallback, useEffect, useState } from 'react'
import { useHistory } from 'react-router-dom'
import { Input } from '../components/Input'
import { Message } from '../components/Message'
import { Connectivity, useConnectivity } from '../hooks/useConnectivity'
import { useSwap } from '../contexts/swap'
import { routes } from '../routes'

type Direction = 'SCtoSF' | 'SFtoSC'

export function CreateNewSwap() {
  const history = useHistory()
  const connectivity = useConnectivity()
  const { transaction, setTransaction, clearTransaction } = useSwap()

  useEffect(() => {
    clearTransaction()
  }, [])

  const [direction, setDirection] = useState<Direction>('SCtoSF')
  const [sc, setSc] = useState<string>()
  const [sf, setSf] = useState<string>()
  const offerSc = direction === 'SCtoSF'

  const handleCreate = useCallback(
    (sc: number, sf: number) => {
      const func = async () => {
        const offer = offerSc ? `${sc}SC` : `${sf}SF`
        const receive = offerSc ? `${sf}SF` : `${sc}SC`

        try {
          const response = await axios({
            method: 'post',
            url: 'http://localhost:8080/api/create',
            headers: {
              'Content-Type': 'application/json',
            },
            data: {
              offer,
              receive,
            },
          })

          setTransaction(response.data.raw)
          history.push(routes.waitingForCounterpartyToAccept)
        } catch (e) {
          if (e instanceof Error) {
            console.log(e.message)
          }
          triggerErrorToast('Error creating swap transaction')
        }
      }
      func()
    },
    [offerSc, history, setTransaction]
  )

  const isValues = sc && sf
  const connError = getConnError(connectivity)

  return (
    <Flex direction="column" align="center" gap="3">
      <Flex direction="column" align="center" css={{ width: '100%' }}>
        <Box css={{ width: '100%', order: offerSc ? 1 : 3 }}>
          <Input
            currency="SC"
            type="decimal"
            disabled={!!transaction}
            value={sc}
            onChange={setSc}
            isOffer={offerSc}
          />
        </Box>
        <Box css={{ width: '100%', order: offerSc ? 3 : 1 }}>
          <Input
            currency="SF"
            type="integer"
            disabled={!!transaction}
            value={sf}
            onChange={setSf}
            isOffer={!offerSc}
          />
        </Box>
        <Box css={{ height: '$2', zIndex: 1, order: 2 }}>
          <Box
            onClick={() => {
              if (transaction) {
                return
              }
              setDirection(direction === 'SCtoSF' ? 'SFtoSC' : 'SCtoSF')
            }}
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
                backgroundColor: '$gray7',
                borderRadius: '$4',
                position: 'absolute',
                transform: 'translate(-50%, -50%)',
                left: '50%',
                color: '$hiContrast',
                top: '50%',
                height: '30px',
                width: '30px',
                transition: 'background 0.1s linear',
                '&:hover': !transaction && {
                  backgroundColor: '$primary10',
                },
              }}
            >
              <ArrowDownIcon />
            </Flex>
          </Box>
        </Box>
      </Flex>
      <Message
        message={`
          The party that contributes SC is responsible for paying the miner
          fee.
      `}
      />
      {connError && <Message variant="red" message={connError} />}
      <Button
        size="3"
        disabled={!connectivity.all || !isValues}
        variant="green"
        css={{ width: '100%', textAlign: 'center' }}
        onClick={() => handleCreate(Number(sc), Number(sf))}
      >
        Generate swap
      </Button>
    </Flex>
  )
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
