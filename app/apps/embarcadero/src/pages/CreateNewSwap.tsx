import {
  ArrowDown16,
  Box,
  Button,
  Flex,
  triggerErrorToast,
} from '@siafoundation/design-system'
import axios from 'axios'
import { useCallback, useEffect, useState } from 'react'
import { Input } from '../components/Input'
import { Message } from '../components/Message'
import { useConnectivity } from '../hooks/useConnectivity'
import { useSwap } from '../contexts/swap'
import { api } from '../config'
import BigNumber from 'bignumber.js'
import { ErrorMessageConn } from '../components/ErrorMessageConn'

type Direction = 'SCtoSF' | 'SFtoSC'

export function CreateNewSwap() {
  const connectivity = useConnectivity()
  const { txn, loadTxn, resetTxn } = useSwap()

  useEffect(() => {
    resetTxn()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const [direction, setDirection] = useState<Direction>('SCtoSF')
  const [sc, setSc] = useState<BigNumber>()
  const [sf, setSf] = useState<BigNumber>()
  const offerSc = direction === 'SCtoSF'

  const isValues = sc && sf
  const readyToCreate = connectivity.all && isValues

  const handleCreate = useCallback(() => {
    if (!readyToCreate) {
      return
    }

    const func = async () => {
      const offer = offerSc ? `${sc.toString()}SC` : `${sf.toString()}SF`
      const receive = offerSc ? `${sf.toString()}SF` : `${sc.toString()}SC`

      try {
        const response = await axios({
          method: 'post',
          url: `${api}/api/create`,
          headers: {
            'Content-Type': 'application/json',
          },
          data: {
            offer,
            receive,
          },
        })

        loadTxn(response.data.swap)
      } catch (e) {
        if (e instanceof Error) {
          console.log(e.message)
        }
        triggerErrorToast('Error creating swap transaction')
      }
    }
    func()
  }, [sc, sf, readyToCreate, offerSc, loadTxn])

  return (
    <Flex direction="column" align="center" gap="3">
      <Flex direction="column" align="center" css={{ width: '100%' }}>
        <Box css={{ width: '100%', order: offerSc ? 1 : 3 }}>
          <Input
            currency="SC"
            type="decimal"
            disabled={!!txn}
            value={sc}
            onChange={setSc}
            isOffer={offerSc}
          />
        </Box>
        <Box css={{ width: '100%', order: offerSc ? 3 : 1 }}>
          <Input
            currency="SF"
            type="integer"
            disabled={!!txn}
            value={sf}
            onChange={setSf}
            isOffer={!offerSc}
          />
        </Box>
        <Box css={{ height: '$2', zIndex: 1, order: 2 }}>
          <Box
            onClick={() => {
              if (txn) {
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
                '&:hover': !txn && {
                  backgroundColor: '$primary10',
                },
              }}
            >
              <ArrowDown16 />
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
      <ErrorMessageConn />
      <Button
        size="3"
        disabled={!readyToCreate}
        variant="green"
        css={{ width: '100%', textAlign: 'center' }}
        onClick={() => handleCreate()}
      >
        Generate swap
      </Button>
    </Flex>
  )
}
