import {
  Toaster,
  Background,
  ScrollArea,
  Dialog,
  DialogContent,
  Button,
} from '@siafoundation/design-system'
import React from 'react'
import { Footer } from './Footer'
import { Navbar } from './Navbar'
import { SwapLayout } from '../SwapLayout'
import { useDialog } from '../../contexts/dialog'
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
