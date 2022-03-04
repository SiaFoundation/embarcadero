import Joi from 'joi'

const publicKey = Joi.object({
  algorithm: Joi.string().required(),
  key: Joi.string().required(),
})

const unlockConditions = Joi.object({
  timelock: Joi.number().required(),
  publickeys: Joi.array().items(publicKey).required(),
  signaturesrequired: Joi.number().required(),
})

const siacoinInput = Joi.object({
  parentid: Joi.string().required(),
  unlockconditions: unlockConditions.required(),
})

const siacoinOutput = Joi.object({
  value: Joi.string().required(),
  unlockhash: Joi.string().required(),
})

const siafundInput = Joi.object({
  parentid: Joi.string().required(),
  unlockconditions: unlockConditions.required(),
  claimunlockhash: Joi.string().required(),
})

const siafundOutput = Joi.object({
  value: Joi.string().required(),
  unlockhash: Joi.string().required(),
  claimstart: Joi.string().required(),
})

const signature = Joi.object({
  parentid: Joi.string().required(),
  publickeyindex: Joi.number().required(),
  timelock: Joi.number().required(),
  coveredfields: Joi.object().required(),
  signature: Joi.string().required(),
})

export const swapTxnSchema = Joi.object({
  siacoinInputs: Joi.alternatives(Joi.array().items(siacoinInput), null),
  siacoinOutputs: Joi.alternatives(Joi.array().items(siacoinOutput), null),
  siafundInputs: Joi.alternatives(Joi.array().items(siafundInput), null),
  siafundOutputs: Joi.alternatives(Joi.array().items(siafundOutput), null),
  signatures: Joi.alternatives(Joi.array().items(signature), null),
})
