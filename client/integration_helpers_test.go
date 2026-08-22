//go:build integration

package client_test

import (
	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/accreditations"
	"github.com/TheWFA/api-client-go/clubs"
	"github.com/TheWFA/api-client-go/competitions"
	"github.com/TheWFA/api-client-go/locations"
	"github.com/TheWFA/api-client-go/organisations"
	"github.com/TheWFA/api-client-go/persons"
	"github.com/TheWFA/api-client-go/seasons"
	"github.com/TheWFA/api-client-go/suspensions"
	"github.com/TheWFA/api-client-go/teams"
	"github.com/TheWFA/api-client-go/ties"
)

func teamsListQuery(n int) teams.ListQuery {
	return teams.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func playersQuery() teams.PlayersQuery {
	return teams.PlayersQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(5)}}
}

func staffQuery() teams.StaffQuery {
	return teams.StaffQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(5)}}
}

func registrationsQuery() teams.RegistrationsQuery {
	return teams.RegistrationsQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(5)}}
}

func clubsListQuery(n int) clubs.ListQuery {
	return clubs.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func competitionsListQuery(n int) competitions.ListQuery {
	return competitions.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func organisationsListQuery(n int) organisations.ListQuery {
	return organisations.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func seasonsListQuery(n int) seasons.ListQuery {
	return seasons.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func accreditationsListQuery(n int) accreditations.ListQuery {
	return accreditations.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func personsListQuery(n int) persons.ListQuery {
	return persons.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func personSuspensionsQuery() suspensions.PersonQuery {
	return suspensions.PersonQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(5)}}
}

func locationsListQuery(n int) locations.ListQuery {
	return locations.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func suspensionsListQuery(n int) suspensions.ListQuery {
	return suspensions.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}

func tiesListQuery(n int) ties.ListQuery {
	return ties.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(n)}}
}
