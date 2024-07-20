const Spinner = ({ show = false, style = {} }) =>
    show ? (
        <span style={{ ...style, fontSize: '14px' }}>
            <i className='fa fa-circle-notch fa-spin' style={{ color: '#0DADEA' }} />
        </span>
    ) : null

export default Spinner