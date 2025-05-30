<script lang="ts" setup>
import user from '@/api/panel/user'
import bgImg from '@/assets/images/login_bg.webp'
import logoImg from '@/assets/images/logo.png'
import { addDynamicRoutes } from '@/router'
import { useThemeStore, useUserStore } from '@/store'
import { getLocal, removeLocal, setLocal } from '@/utils'
import { rsaEncrypt } from '@/utils/encrypt'
import { useGettext } from 'vue3-gettext'

const { $gettext } = useGettext()
const router = useRouter()
const route = useRoute()
const query = route.query
const { data: key, loading: isLoading } = useRequest(user.key, { initialData: '' })
const { data: isLogin } = useRequest(user.isLogin, { initialData: false })

interface LoginInfo {
  username: string
  password: string
  safe_login: boolean
  pass_code: string
}

const loginInfo = ref<LoginInfo>({
  username: '',
  password: '',
  safe_login: true,
  pass_code: ''
})

const localLoginInfo = getLocal('loginInfo') as LoginInfo
if (localLoginInfo) {
  loginInfo.value.username = localLoginInfo.username || ''
  loginInfo.value.password = localLoginInfo.password || ''
}

const userStore = useUserStore()
const themeStore = useThemeStore()
const loging = ref<boolean>(false)
const isRemember = useStorage('isRemember', false)
const showTwoFA = ref(false)

const logo = computed(() => themeStore.logo || logoImg)

async function handleLogin() {
  const { username, password, pass_code, safe_login } = loginInfo.value
  if (!username || !password) {
    window.$message.warning($gettext('Please enter username and password'))
    return
  }
  if (!key) {
    window.$message.warning(
      $gettext('Failed to get encryption public key, please refresh the page and try again')
    )
    return
  }
  useRequest(
    user.login(
      rsaEncrypt(username, String(unref(key))),
      rsaEncrypt(password, String(unref(key))),
      pass_code,
      safe_login
    )
  ).onSuccess(async () => {
    loging.value = true
    window.$notification?.success({ title: $gettext('Login successful!'), duration: 2500 })
    if (isRemember.value) {
      setLocal('loginInfo', { username, password })
    } else {
      removeLocal('loginInfo')
    }

    await addDynamicRoutes()
    useRequest(user.info()).onSuccess(({ data }) => {
      userStore.set(data as any)
    })
    if (query.redirect) {
      const path = query.redirect as string
      Reflect.deleteProperty(query, 'redirect')
      await router.push({ path, query })
    } else {
      await router.push('/')
    }
  })
  loging.value = false
}

const isTwoFA = () => {
  const { username } = loginInfo.value
  if (!username) {
    return
  }
  useRequest(user.isTwoFA(username))
    .onSuccess(({ data }) => {
      showTwoFA.value = Boolean(data)
    })
    .onError(() => {
      showTwoFA.value = false
    })
}

watch(isLogin, async () => {
  if (isLogin) {
    await addDynamicRoutes()
    useRequest(user.info()).onSuccess(({ data }) => {
      userStore.set(data as any)
    })
    if (query.redirect) {
      const path = query.redirect as string
      Reflect.deleteProperty(query, 'redirect')
      await router.push({ path, query })
    } else {
      await router.push('/')
    }
  }
})
</script>

<template>
  <AppPage :show-footer="true" :style="{ backgroundImage: `url(${bgImg})` }" bg-cover>
    <div m-auto min-w-345 f-c-c bg-white bg-opacity-60 p-15 card-shadow dark:bg-dark>
      <div w-480 flex-col px-20 py-35>
        <h5 color="#6a6a6a" f-c-c text-24 font-normal>
          <n-image :src="logo" height="50" preview-disabled mr-10 />{{ themeStore.name }}
        </h5>
        <div mt-30>
          <n-input
            v-model:value="loginInfo.username"
            :maxlength="32"
            autofocus
            class="h-50 items-center pl-10 text-16"
            :placeholder="$gettext('Username')"
            :on-blur="isTwoFA"
          />
        </div>
        <div mt-30>
          <n-input
            v-model:value="loginInfo.password"
            :maxlength="32"
            class="h-50 items-center pl-10 text-16"
            :placeholder="$gettext('Password')"
            type="password"
            show-password-on="click"
            @keydown.enter="handleLogin"
          />
        </div>
        <div v-if="showTwoFA" mt-30>
          <n-input
            v-model:value="loginInfo.pass_code"
            :maxlength="6"
            class="h-50 items-center pl-10 text-16"
            :placeholder="$gettext('2FA Code')"
            type="text"
            @keydown.enter="handleLogin"
          />
        </div>

        <div mt-20>
          <n-flex>
            <n-checkbox v-model:checked="loginInfo.safe_login" :label="$gettext('Safe Login')" />
            <n-checkbox v-model:checked="isRemember" :label="$gettext('Remember Me')" />
          </n-flex>
        </div>

        <div mt-20>
          <n-button
            :loading="isLoading || loging"
            :disabled="isLoading || loging"
            type="primary"
            h-50
            w-full
            text-16
            @click="handleLogin"
          >
            {{ $gettext('Login') }}
          </n-button>
        </div>
      </div>
    </div>
  </AppPage>
</template>
