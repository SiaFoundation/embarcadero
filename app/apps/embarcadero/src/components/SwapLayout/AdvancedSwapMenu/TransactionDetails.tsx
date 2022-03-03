import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
  Box,
  Codeblock,
  Text,
} from '@siafoundation/design-system'
import { SwapTransaction, useSwap } from '../../../contexts/swap'

export function TransactionDetails() {
  const { transaction } = useSwap()

  if (!transaction) {
    return (
      <Text css={{ p: '$3 $2', color: '$slate9' }}>
        Load a swap to view details
      </Text>
    )
  }

  const keys = Object.keys(transaction) as (keyof SwapTransaction)[]

  return (
    <Box css={{ pl: '$3' }}>
      <Accordion type="single">
        {keys.map((key) => (
          <AccordionItem key={key} value={key}>
            <AccordionTrigger>
              <Text size="3" css={{ fontWeight: 500 }}>
                {key} ({transaction[key]?.length || 0})
              </Text>
            </AccordionTrigger>
            <AccordionContent>
              <pre>
                <Codeblock css={{ overflow: 'auto' }}>
                  {JSON.stringify(transaction[key], null, 2)}
                </Codeblock>
              </pre>
            </AccordionContent>
          </AccordionItem>
        ))}
      </Accordion>
    </Box>
  )
}
