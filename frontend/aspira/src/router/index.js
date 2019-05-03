import Login from '@/components/auth/Login'
import Home from '@/components/process/Home'
import Search from '@/components/process/Search'
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
  path: '/search',
  name: 'Search',
  component: Search
},
{
  path: '/claims',
  name: 'Claims',
  component: Claims
}
]

export default routes
