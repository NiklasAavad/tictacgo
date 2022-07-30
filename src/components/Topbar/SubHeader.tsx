import './SubHeader.css'

export const SubHeader: React.FC = () => {
    const lobby = <a href="/lobby">Lobby</a>
    const user = <a href="/user">User</a>
    const help = <a href="/help">Help</a>

    return <div className='sub-header'>{lobby} | {user} | {help}</div>
}