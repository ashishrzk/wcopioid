import Login from '@/components/auth/Login'
import Home from '@/components/process/Home'
import Claims from '@/components/process/Claims'

const routes = [{
  path: '/home',
  name: 'Home',
  component: Home
},
{
  path: '/',
  name: 'Login',
  component: Login
},
{
  path: '/claims',
  name: 'Claims',
  component: Claims
}
]

export default routes
