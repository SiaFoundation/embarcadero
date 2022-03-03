import { Toaster, Background, ScrollArea } from '@siafoundation/design-system'
import React from 'react'
import { Footer } from './Footer'
import { Navbar } from './Navbar'
import { SwapLayout } from '../SwapLayout'
import { Dialogs } from './Dialogs'

type Props = {
  children: React.ReactNode
}

export function Layout({ children }: Props) {
  return (
    <ScrollArea>
      <Dialogs />
      <Toaster />
      <Background level="1" />
      <Navbar />
      <SwapLayout>{children}</SwapLayout>
      <Footer />
    </ScrollArea>
  )
}
