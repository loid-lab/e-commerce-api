import { z } from "zod";

export const paymentSchema = z.object({
    orderID: z.number().int(),
    method: z.string().min(1),
    provider: z.string().min(1),
    refID: z.string().min(1),
    status: z.enum(["paid", "pending", "failed"])
})

export type PaymentSchema = z.infer<typeof paymentSchema>

