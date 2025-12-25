package kafkastorage

func (p *Publisher) Close() error {
	return p.w.Close()
}


