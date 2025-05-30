import { z } from "zod";

export const orderSchema = z.object({
    shippingAddressID: z.number().int()
    // additional fields if needed
})

export type OrderSchema = z.infer<typeof orderSchema>