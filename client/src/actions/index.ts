import { ActionError, defineAction } from 'astro:actions'
import { z } from 'astro:schema'
import axios, { AxiosError } from 'axios'

export const server = {
  login: defineAction({
    accept: 'form',
    input: z.object({
      username: z.string(),
      password: z.string()
    }),
    handler: async ({ username, password }) => {
      try {
        console.log(username, password)
        const res = await axios.post('http://localhost:8080/login', {
          username: username,
          password: password
        })

        return { token: res.data.token }
      } catch (err) {
        console.log((err as AxiosError).response?.data)
        throw new ActionError({
          code: 'UNAUTHORIZED',
          message: 'Incorrect Credentials'
        })
      }
    }
  })
}
