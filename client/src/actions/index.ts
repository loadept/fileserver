import { ActionError, defineAction } from 'astro:actions'
import { z } from 'astro:schema'
import { AxiosError } from 'axios'
import apiClient from '../config/apiClient'

export const server = {
  login: defineAction({
    accept: 'form',
    input: z.object({
      username: z.string(),
      password: z.string()
    }),
    handler: async ({ username, password }) => {
      try {
        const res = await apiClient.post('/login', {
          username: username,
          password: password
        })

        return { token: res.data.token }
      } catch (err) {
        if (err instanceof AxiosError) {
          console.log('Error:', (err as AxiosError).response?.data)
        } else {
          console.log('Error:', err)
        }
        throw new ActionError({
          code: 'UNAUTHORIZED',
          message: 'Incorrect Credentials'
        })
      }
    }
  })
}
