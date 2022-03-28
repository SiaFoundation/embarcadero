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
import { ToggleInputs } from '../components/ToggleInputs'

type Direction = 'SCtoSF' | 'SFtoSC'

export function CreateNewSwap() {
  const connectivity = useConnectivity()
  const { txn, loadTxn, resetTxn } = useSwap()

  useEffect(() => {
    resetTxn()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const [direction, setDirection] = useState<Direction>('SFtoSC')
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
            tabIndex={offerSc ? 1 : 3}
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
            tabIndex={offerSc ? 3 : 1}
            disabled={!!txn}
            value={sf}
            onChange={setSf}
            isOffer={!offerSc}
          />
        </Box>
        <ToggleInputs
          onToggle={() => {
            if (txn) {
              return
            }
            setDirection(direction === 'SCtoSF' ? 'SFtoSC' : 'SCtoSF')
          }}
          disabled={!txn}
        />
      </Flex>
      <Flex direction="column" align="center" gap="1-5">
        {direction === 'SCtoSF' && (
          <Message
            variant="info"
            message={`
            The party sending SC is responsible for paying the miner fee.
          `}
          />
        )}
        {direction === 'SFtoSC' && (
          <Message
            variant="info"
            message={`
            The party sending SF will receive a separate SC payout for with file contract dividends.
          `}
          />
        )}
        <ErrorMessageConn />
        <Button
          size="3"
          disabled={!readyToCreate}
          variant="accent"
          css={{ width: '100%', textAlign: 'center', borderRadius: '$2' }}
          onClick={() => handleCreate()}
        >
          Generate swap
        </Button>
      </Flex>
    </Flex>
  )
}
