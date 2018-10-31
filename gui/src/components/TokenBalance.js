import React from 'react'
import PropTypes from 'prop-types'
import numeral from 'numeral'
import { BigNumber } from 'bignumber.js'
import { Tooltip } from '@material-ui/core'
import Typography from '@material-ui/core/Typography'
import PaddedCard from 'components/PaddedCard'

const WEI_PER_TOKEN = new BigNumber(10 ** 18)

const formatBalance = (val) => {
  const b = new BigNumber(val)
  const tokenBalance = b.dividedBy(WEI_PER_TOKEN).toNumber()
  return {formatted: numeral(tokenBalance).format('0.200000a'), unformatted: tokenBalance}
}

const valAndTooltip = ({value, title, error}) => {
  let val, tooltip

  if (error) {
    val = error
    tooltip = 'Error'
  } else if (value == null) {
    val = '...'
    tooltip = 'Loading...'
  } else {
    const balance = formatBalance(value)
    val = balance.formatted
    tooltip = balance.unformatted
  }

  return {val, tooltip}
}

const TokenBalance = props => {
  const {val, tooltip} = valAndTooltip(props)

  return (
    <PaddedCard>
      <Typography gutterBottom variant='headline' component='h2'>
        {props.title}
      </Typography>
      <Typography variant='body1' color='textSecondary'>
        <Tooltip title={tooltip} placement='left'>
          <span>{val}</span>
        </Tooltip>
      </Typography>
    </PaddedCard>
  )
}

TokenBalance.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number
  ]),
  className: PropTypes.string,
  error: PropTypes.string
}

export default TokenBalance
