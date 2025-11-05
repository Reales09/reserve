import 'package:get/get.dart';
import 'package:rupu/domain/datasource/horizontal_properties_datasource.dart';
import 'package:rupu/domain/entities/horizontal_properties_page.dart';
import 'package:rupu/domain/entities/horizontal_property_action_result.dart';
import 'package:rupu/domain/entities/horizontal_property_create_result.dart';
import 'package:rupu/domain/entities/horizontal_property_detail.dart';
import 'package:rupu/domain/entities/horizontal_property_residents_page.dart';
import 'package:rupu/domain/entities/horizontal_property_unit_detail.dart';
import 'package:rupu/domain/entities/horizontal_property_units_page.dart';
import 'package:rupu/domain/entities/horizontal_property_update_result.dart';
import 'package:rupu/domain/entities/horizontal_property_voting_groups.dart';
import 'package:rupu/domain/infrastructure/datasources/horizontal_properties_datasource_impl.dart';
import 'package:rupu/domain/entities/voting_results_result.dart';
import 'package:rupu/domain/entities/votings_result.dart';
import 'package:rupu/domain/infrastructure/mappers/horizontal_properties_mapper.dart';
import 'package:rupu/domain/infrastructure/mappers/voting_result_mapper.dart';
import 'package:rupu/domain/infrastructure/mappers/voting_mapper.dart';
import 'package:rupu/domain/repositories/horizontal_properties_repository.dart';
import 'package:rupu/presentation/views/login/login_controller.dart';

class HorizontalPropertiesRepositoryImpl
    extends HorizontalPropertiesRepository {
  final HorizontalPropertiesDatasource datasource;

  HorizontalPropertiesRepositoryImpl({HorizontalPropertiesDatasource? datasource})
      : datasource =
            datasource ?? HorizontalPropertiesDatasourceImpl();

  @override
  Future<HorizontalPropertiesPage> getHorizontalProperties({
    Map<String, dynamic>? query,
  }) async {
    final response = await datasource.getHorizontalProperties(
      query: _withBusinessQuery(query),
    );
    return HorizontalPropertiesMapper.responseToEntity(response);
  }

  @override
  Future<HorizontalPropertyCreateResult> createHorizontalProperty({
    required Map<String, dynamic> data,
  }) async {
    final response = await datasource.createHorizontalProperty(
      data: data,
      query: _withBusinessQuery(),
    );
    return HorizontalPropertiesMapper.detailResponseToCreateResult(response);
  }

  @override
  Future<HorizontalPropertyActionResult> deleteHorizontalProperty({
    required int id,
  }) async {
    final response = await datasource.deleteHorizontalProperty(
      id: id,
      query: _withBusinessQuery(),
    );
    return HorizontalPropertiesMapper.simpleResponseToActionResult(response);
  }

  @override
  Future<HorizontalPropertyDetail?> getHorizontalPropertyDetail({
    required int id,
  }) async {
    final response = await datasource.getHorizontalPropertyDetail(
      id: id,
      query: _withBusinessQuery(),
    );
    return HorizontalPropertiesMapper.detailResponseToDetail(response);
  }

  @override
  Future<HorizontalPropertyUpdateResult> updateHorizontalProperty({
    required int id,
    required Map<String, dynamic> data,
    String? logoFilePath,
    String? logoFileName,
    String? navbarImagePath,
    String? navbarImageFileName,
  }) async {
    final response = await datasource.updateHorizontalProperty(
      id: id,
      data: data,
      logoFilePath: logoFilePath,
      logoFileName: logoFileName,
      navbarImagePath: navbarImagePath,
      navbarImageFileName: navbarImageFileName,
      query: _withBusinessQuery(),
    );
    return HorizontalPropertiesMapper.detailResponseToUpdateResult(response);
  }

  @override
  Future<HorizontalPropertyUnitsPage> getHorizontalPropertyUnits({
    required int id,
    Map<String, dynamic>? query,
  }) async {
    final response = await datasource.getHorizontalPropertyUnits(
      id: id,
      query: _withBusinessQuery(query),
    );
    return HorizontalPropertiesMapper.unitsResponseToEntity(response);
  }

  @override
  Future<HorizontalPropertyUnitDetailResult> getHorizontalPropertyUnitDetail({
    required int unitId,
  }) async {
    try {
      final response = await datasource.getHorizontalPropertyUnitDetail(
        unitId: unitId,
        query: _withBusinessQuery(),
      );
      return HorizontalPropertiesMapper.unitDetailResponseToEntity(response);
    } catch (_) {
      return const HorizontalPropertyUnitDetailResult(
        success: false,
        message: 'No se pudo obtener el detalle de la unidad.',
      );
    }
  }

  @override
  Stream<LiveVotingResult> getLiveVotingStream({
    required int propertyId,
    required int groupId,
    required int votingId,
  }) {
    return datasource
        .getLiveVotingStream(
      propertyId: propertyId,
      groupId: groupId,
      votingId: votingId,
      query: _withBusinessQuery(),
    )
        .map((response) {
      return LiveVotingMapper.fromJson(response.data);
    });
  }

  @override
  Future<VotingResultsResult> getVotingResults({
    required int propertyId,
    required int groupId,
    required int votingId,
  }) async {
    try {
      final response = await datasource.getVotingResults(
        propertyId: propertyId,
        groupId: groupId,
        votingId: votingId,
        query: _withBusinessQuery(),
      );
      return VotingResultMapper.fromJson(response.data);
    } catch (_) {
      return const VotingResultsResult(
        success: false,
        message: 'No se pudieron cargar los resultados de la votación.',
        data: [],
      );
    }
  }

  @override
  Future<VotingOptionsResult> getVotingOptions({
    required int propertyId,
    required int groupId,
    required int votingId,
  }) async {
    try {
      final response = await datasource.getVotingOptions(
        propertyId: propertyId,
        groupId: groupId,
        votingId: votingId,
        query: _withBusinessQuery(),
      );
      return VotingOptionMapper.fromJson(response.data);
    } catch (_) {
      return const VotingOptionsResult(
        success: false,
        message: 'No se pudieron cargar las opciones de votación.',
        data: [],
      );
    }
  }

  @override
  Future<bool> createVotingOption({
    required int propertyId,
    required int groupId,
    required int votingId,
    required Map<String, dynamic> data,
  }) async {
    try {
      await datasource.createVotingOption(
        propertyId: propertyId,
        groupId: groupId,
        votingId: votingId,
        data: data,
        query: _withBusinessQuery(),
      );
      return true;
    } catch (_) {
      return false;
    }
  }

  @override
  Future<HorizontalPropertyActionResult> deleteVotingOption({
    required int propertyId,
    required int groupId,
    required int votingId,
    required int optionId,
  }) async {
    try {
      final response = await datasource.deleteVotingOption(
        propertyId: propertyId,
        groupId: groupId,
        votingId: votingId,
        optionId: optionId,
        query: _withBusinessQuery(),
      );
      return HorizontalPropertiesMapper.simpleResponseToActionResult(response);
    } catch (_) {
      return const HorizontalPropertyActionResult(
        success: false,
        message: 'No se pudo eliminar la opción de votación.',
      );
    }
  }

  @override
  Future<HorizontalPropertyUnitDetailResult> createHorizontalPropertyUnit({
    required int propertyId,
    required Map<String, dynamic> data,
  }) async {
    try {
      final response = await datasource.createHorizontalPropertyUnit(
        data: data,
        query: _withBusinessQuery(
          propertyId > 0 ? {'business_id': propertyId} : null,
        ),
      );
      return HorizontalPropertiesMapper.unitDetailResponseToEntity(response);
    } catch (_) {
      return const HorizontalPropertyUnitDetailResult(
        success: false,
        message: 'No se pudo crear la unidad.',
      );
    }
  }

  @override
  Future<HorizontalPropertyUnitDetailResult> updateHorizontalPropertyUnit({
    required int propertyId,
    required int unitId,
    required Map<String, dynamic> data,
  }) async {
    try {
      final response = await datasource.updateHorizontalPropertyUnit(
        unitId: unitId,
        data: data,
        query: _withBusinessQuery(
          propertyId > 0 ? {'business_id': propertyId} : null,
        ),
      );
      return HorizontalPropertiesMapper.unitDetailResponseToEntity(response);
    } catch (_) {
      return const HorizontalPropertyUnitDetailResult(
        success: false,
        message: 'No se pudo actualizar la unidad.',
      );
    }
  }

  @override
  Future<HorizontalPropertyActionResult> deleteHorizontalPropertyUnit({
    required int propertyId,
    required int unitId,
  }) async {
    try {
      final response = await datasource.deleteHorizontalPropertyUnit(
        unitId: unitId,
        query: _withBusinessQuery(
          propertyId > 0 ? {'business_id': propertyId} : null,
        ),
      );
      return HorizontalPropertiesMapper.simpleResponseToActionResult(response);
    } catch (_) {
      return const HorizontalPropertyActionResult(
        success: false,
        message: 'No se pudo eliminar la unidad.',
      );
    }
  }

  @override
  Future<HorizontalPropertyResidentsPage> getHorizontalPropertyResidents({
    required int id,
    Map<String, dynamic>? query,
  }) async {
    final response = await datasource.getHorizontalPropertyResidents(
      id: id,
      query: _withBusinessQuery(query),
    );
    return HorizontalPropertiesMapper.residentsResponseToEntity(response);
  }

  @override
  Future<HorizontalPropertyVotingGroupsResult> getHorizontalPropertyVotingGroups({
    required int id,
  }) async {
    final response = await datasource.getHorizontalPropertyVotingGroups(
      id: id,
      query: _withBusinessQuery({'business_id': id}),
    );
    return HorizontalPropertiesMapper.votingGroupsResponseToEntity(response);
  }

  @override
  Future<bool> createHorizontalPropertyVotingGroup({
    required int propertyId,
    required Map<String, dynamic> data,
  }) async {
    try {
      await datasource.createHorizontalPropertyVotingGroup(
        propertyId: propertyId,
        data: data,
        query: _withBusinessQuery(),
      );
      return true;
    } catch (_) {
      return false;
    }
  }

  @override
  Future<bool> updateHorizontalPropertyVotingGroup({
    required int propertyId,
    required int groupId,
    required Map<String, dynamic> data,
  }) async {
    try {
      await datasource.updateHorizontalPropertyVotingGroup(
        propertyId: propertyId,
        groupId: groupId,
        data: data,
        query: _withBusinessQuery(),
      );
      return true;
    } catch (_) {
      return false;
    }
  }

  @override
  Future<HorizontalPropertyActionResult> deleteHorizontalPropertyVotingGroup({
    required int propertyId,
    required int groupId,
  }) async {
    try {
      final response = await datasource.deleteHorizontalPropertyVotingGroup(
        propertyId: propertyId,
        groupId: groupId,
        query: _withBusinessQuery(),
      );
      return HorizontalPropertiesMapper.simpleResponseToActionResult(response);
    } catch (_) {
      return const HorizontalPropertyActionResult(
        success: false,
        message: 'No se pudo eliminar el grupo de votación.',
      );
    }
  }

  Map<String, dynamic>? _withBusinessQuery([Map<String, dynamic>? query]) {
    final loginController =
        Get.isRegistered<LoginController>() ? Get.find<LoginController>() : null;

    final baseQuery = query == null
        ? <String, dynamic>{}
        : Map<String, dynamic>.from(query);

    final existingBusinessId = baseQuery['business_id'];
    if (existingBusinessId != null) {
      if (existingBusinessId is num && existingBusinessId > 0) {
        return baseQuery;
      }
      if (existingBusinessId is String) {
        final trimmed = existingBusinessId.trim();
        if (trimmed.isNotEmpty && trimmed != '0') {
          return baseQuery;
        }
      }
      baseQuery.remove('business_id');
    }

    final businessId = loginController?.selectedBusinessId;

    if (businessId == null) {
      return baseQuery.isEmpty ? null : baseQuery;
    }

    baseQuery['business_id'] = businessId;
    return baseQuery;
  }

  @override
  Future<VotingsResult> getVotings({
    required int propertyId,
    required int groupId,
  }) async {
    try {
      final response = await datasource.getVotings(
        propertyId: propertyId,
        groupId: groupId,
        query: _withBusinessQuery(),
      );
      return VotingMapper.fromJson(response.data);
    } catch (_) {
      return const VotingsResult(
        success: false,
        message: 'No se pudieron cargar las votaciones.',
        data: [],
      );
    }
  }

  @override
  Future<HorizontalPropertyActionResult> activateVoting({
    required int propertyId,
    required int groupId,
    required int votingId,
  }) async {
    try {
      final response = await datasource.activateVoting(
        propertyId: propertyId,
        groupId: groupId,
        votingId: votingId,
        query: _withBusinessQuery(),
      );
      return HorizontalPropertiesMapper.simpleResponseToActionResult(response);
    } catch (_) {
      return const HorizontalPropertyActionResult(
        success: false,
        message: 'No se pudo activar la votación.',
      );
    }
  }

  @override
  Future<HorizontalPropertyActionResult> deactivateVoting({
    required int propertyId,
    required int groupId,
    required int votingId,
  }) async {
    try {
      final response = await datasource.deactivateVoting(
        propertyId: propertyId,
        groupId: groupId,
        votingId: votingId,
        query: _withBusinessQuery(),
      );
      return HorizontalPropertiesMapper.simpleResponseToActionResult(response);
    } catch (_) {
      return const HorizontalPropertyActionResult(
        success: false,
        message: 'No se pudo desactivar la votación.',
      );
    }
  }
}
