import { Header } from './Header'
import { SubHeader } from './SubHeader'
import './Topbar.css'

export const Topbar: React.FC = () => {
    return <div className='topbar'>
        <Header/>
        <SubHeader/>
    </div>
}