import { z } from "zod"

export const authSchema = z.object({
    email: z.string()
    .email().
    regex(/^[^\s@]+@[^\s@]+\.[^\s@]+$/, { message: "Invalid email structure" }),
    password: z.string()
    .min(8, { message: "Password must be at least 8 characters long" })
    .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/, { message: "Invalid password" })
})

export type AuthSchema = z.infer<typeof authSchema>;