package mongodb

import (
	"log"

	"github.com/google/uuid"
	"github.com/looplab/eventhorizon/commandhandler/bus"
	"github.com/nvloff/f3_payments_service/internal/domain"
	"github.com/nvloff/f3_payments_service/internal/persistence"

	eh "github.com/looplab/eventhorizon"
	eventbus "github.com/looplab/eventhorizon/eventbus/local"
	eventstore "github.com/looplab/eventhorizon/eventstore/mongodb"
	repo "github.com/looplab/eventhorizon/repo/mongodb"
	"github.com/looplab/eventhorizon/repo/version"
)

const defaultMongoHost = "localhost:27017"
const defaultRepoDbPrefix = "form3"
const defaultEventsDbPrefix = "form3"
const defaultReadRepoCollectionName = "payments"

// Config stores mongodb connection options and domain storage config
type Config struct {
	URL            string
	RepoDbPrefix   string
	EventsDbPrefix string
	Collection     string
}

// BuildDomain builds a payment domain with mongodb storage configuration
func BuildDomain(config Config) *domain.PaymentsDomain {
	d := domain.NewPaymentsDomain(initDomainConfig(config))
	d.Setup()

	return d
}

// Initialise the Domain Infrastructure with MongoDB EventStore and Projector
func initDomainConfig(config Config) *domain.Config {
	url := config.URL
	if len(url) == 0 {
		url = defaultMongoHost
	}

	repoDbPrefix := config.RepoDbPrefix
	if len(repoDbPrefix) == 0 {
		repoDbPrefix = defaultRepoDbPrefix
	}

	eventsDbPrefix := config.EventsDbPrefix
	if len(eventsDbPrefix) == 0 {
		eventsDbPrefix = defaultEventsDbPrefix
	}

	collection := config.Collection
	if len(collection) == 0 {
		collection = defaultReadRepoCollectionName
	}

	eventID := uuid.New()

	paymentsRepo := buildVersionRepo(
		buildRepository(url, repoDbPrefix, collection))

	// Setup the domain.
	return &domain.Config{
		EventStore:  buildEventStore(url, eventsDbPrefix),
		EventBus:    buildEventBus(),
		CommandBus:  buildCommandBus(),
		EventID:     eventID,
		PaymentRepo: paymentsRepo,
	}
}

func buildEventStore(url string, dbPrefix string) eh.EventStore {
	// Create the event store.
	eventStore, err := eventstore.NewEventStore(url, dbPrefix)
	if err != nil {
		log.Fatalf("could not create event store: %s", err)
	}

	return eventStore
}

func buildEventBus() eh.EventBus {
	// Create the event bus that distributes events.
	eventBus := eventbus.NewEventBus(nil)
	go func() {
		for e := range eventBus.Errors() {
			log.Printf("eventbus: %s", e.Error())
		}
	}()

	return eventBus
}

func buildCommandBus() *bus.CommandHandler {
	return bus.NewCommandHandler()
}

func buildRepository(url string, dbPrefix string, collection string) eh.ReadWriteRepo {
	// Create the read repositories.
	paymentsRepo, err := repo.NewRepo(url, dbPrefix, collection)
	if err != nil {
		log.Fatalf("could not create payments repository: %s", err)
	}
	paymentsRepo.SetEntityFactory(func() eh.Entity { return &persistence.Payment{} })

	return paymentsRepo
}

func buildVersionRepo(repo eh.ReadWriteRepo) *version.Repo {
	// A version repo is needed for the projector to handle eventual consistency.
	paymentsVersionRepo := version.NewRepo(repo)

	return paymentsVersionRepo
}
