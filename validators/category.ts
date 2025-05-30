import { z } from "zod";

export const categorySchema = z.object({
    name: z.string().min(2),
})

export type CategorySchema = z.infer<typeof categorySchema>