import { z } from "zod";

export const cartSchema = z.object({
    productID: z.number().int(),
    quantity: z.number().int().min(1)
})

export type CartSchema = z.infer<typeof cartSchema>

