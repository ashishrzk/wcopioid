import Login from '@/components/auth/Login'
import Home from '@/components/process/Home'
import Search from '@/components/process/Search'
import Claims from '@/components/process/Claims'
import Claim from '@/components/process/Claim'
import PatientAnalytics from '@/components/process/PatientAnalytics'
import PAnalytics from '@/components/process/PAnalytics'
import Provider from '@/components/process/Provider'

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
},
{
  path: '/panalytics',
  name: 'PAnalytics',
  component: PAnalytics
},
{
  path: '/provideranalytics',
  name: 'Providers',
  component: Provider
},
{
  // path: '/claim'
  path: '/claims/:claim_number',
  name: 'Claim',
  component: Claim
},
{
  path: '/patientanalytics',
  name: 'Patientanalytics',
  component: PatientAnalytics
}
]

export default routes
