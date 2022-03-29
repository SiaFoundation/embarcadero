import { Toaster, AppBackdrop, ScrollArea } from '@siafoundation/design-system'
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
      <AppBackdrop />
      <Navbar />
      <SwapLayout>{children}</SwapLayout>
      <Footer />
    </ScrollArea>
  )
}
