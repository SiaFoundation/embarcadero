import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogTitle,
  DialogTrigger,
} from '@siafoundation/design-system'
import { useDialog } from '../../contexts/dialog'
import { PrivacyDialog } from '../PrivacyDialog'

export function Dialogs() {
  const { dialog, closeDialog } = useDialog()

  return (
    <Dialog open={!!dialog} onOpenChange={() => closeDialog()}>
      <DialogPortal>
        <DialogOverlay />
        <DialogContent>
          {dialog === 'privacy' && <PrivacyDialog />}
        </DialogContent>
      </DialogPortal>
    </Dialog>
  )
}
